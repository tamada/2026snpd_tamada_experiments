import ghidra.app.script.GhidraScript;
import ghidra.app.decompiler.*;
import ghidra.program.model.pcode.PcodeOp;
import ghidra.program.model.pcode.PcodeOpAST;
import ghidra.program.model.pcode.Varnode;
import ghidra.program.model.pcode.HighFunction;
import ghidra.program.model.listing.Function;

import java.io.IOException;
import java.nio.file.Files;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.stream.Collectors;

public class HighPCodeLifter extends GhidraScript {

    @Override
    public void run() throws Exception {
        DecompInterface decompInterface = new DecompInterface();
        decompInterface.openProgram(currentProgram);
        var path = java.nio.file.Path.of(currentProgram.getExecutablePath());

        List<String> jsonOutput = new ArrayList<>();
        jsonOutput.add("{");
        jsonOutput.add(String.format("  \"program\": \"%s\",", currentProgram.getName()));
        jsonOutput.add(String.format("  \"path\": \"%s\",", path));
        jsonOutput.add("  \"functions\": [");

        List<String> functionBlocks = new ArrayList<>();
        Function func = getFirstFunction();

        while (func != null && !monitor.isCancelled()) {
            if (!func.isThunk() && !func.isExternal()) {
                DecompileResults results = decompInterface.decompileFunction(func, 30, monitor);
                if (results != null && results.decompileCompleted()) {
                    functionBlocks.add(getFunctionJson(func, results.getHighFunction()));
                }
            }
            func = getFunctionAfter(func);
        }

        // 関数ブロックをカンマで結合して追加
        jsonOutput.add(String.join(",\n", functionBlocks));
        jsonOutput.add("  ]");
        jsonOutput.add("}");

        outputToFile(currentProgram.getName(), jsonOutput);

        decompInterface.dispose();
    }

    private void outputToFile(String fileName, List<String> outputs) {
        try (var out = Files.newBufferedWriter(java.nio.file.Path.of(fileName + ".json"))) {
            var w = new java.io.PrintWriter(out);
            outputs.stream()
                .forEach(line -> w.println(line));
        } catch(java.io.IOException e) {
            e.printStackTrace();
        }
    }

    private String getFunctionJson(Function func, HighFunction highFunc) {
        List<String> opsJson = new ArrayList<>();
        Iterator<PcodeOpAST> opIter = highFunc.getPcodeOps();

        while (opIter.hasNext()) {
            PcodeOpAST op = opIter.next();
            opsJson.add(getOpJson(op));
        }

        return String.format(
            "    {\n      \"name\": \"%s\",\n      \"ops\": [\n%s\n      ]\n    }",
            func.getName(),
            opsJson.stream().map(s -> "        " + s).collect(Collectors.joining(",\n"))
        );
    }

    private String getOpJson(PcodeOp op) {
        String mnemonic = op.getMnemonic();
        Varnode out = op.getOutput();
        Varnode[] inputs = op.getInputs();

        // 入力Varnodeのリストを作成
        List<String> inputStrings = new ArrayList<>();
        for (Varnode in : inputs) {
            inputStrings.add(String.format("\"%s\"", in.toString()));
        }
        String inputsJson = String.join(", ", inputStrings);

        // 出力Varnodeの有無でフォーマットを分ける
        if (out != null) {
            return String.format(
                "{\"op\": \"%s\", \"out\": \"%s\", \"inputs\": [%s]}",
                mnemonic, out.toString(), inputsJson
            );
        } else {
            return String.format(
                "{\"op\": \"%s\", \"inputs\": [%s]}",
                mnemonic, inputsJson
            );
        }
    }
}